package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type column struct {
	name        string
	typ         string
	constraints string
}

type table struct {
	name    string
	columns []column
}

type schemaState struct {
	tables map[string]*table
}

func main() {
	migrationsDir := flag.String("migrations", filepath.FromSlash("api/migrations"), "Path to migrations directory")
	docPath := flag.String("doc", filepath.FromSlash("docs/ai/database_schema.md"), "Path to schema doc")
	flag.Parse()

	latest, state, files, err := buildSchema(*migrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build schema: %v\n", err)
		os.Exit(1)
	}

	if err := updateDoc(*docPath, latest, state, files); err != nil {
		fmt.Fprintf(os.Stderr, "update doc: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s from %d migration files (latest %s).\n", *docPath, len(files), latest)
}

func buildSchema(migrationsDir string) (string, *schemaState, []string, error) {
	glob := filepath.Join(migrationsDir, "*.up.sql")
	files, err := filepath.Glob(glob)
	if err != nil {
		return "", nil, nil, err
	}
	if len(files) == 0 {
		return "", nil, nil, errors.New("no up migration files found")
	}
	sort.Strings(files)

	state := &schemaState{tables: map[string]*table{}}
	latest := ""

	for _, file := range files {
		base := filepath.Base(file)
		if len(base) >= 4 {
			latest = base[:4]
		}

		content, readErr := os.ReadFile(file)
		if readErr != nil {
			return "", nil, nil, readErr
		}

		for _, stmt := range splitSQLStatements(string(content)) {
			applyStatement(state, stmt)
		}
	}

	return latest, state, files, nil
}

func splitSQLStatements(sql string) []string {
	var out []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	runes := []rune(sql)

	for i, r := range runes {
		switch r {
		case '\'':
			if !inDoubleQuote {
				nextIsQuote := i+1 < len(runes) && runes[i+1] == '\''
				if !nextIsQuote {
					inSingleQuote = !inSingleQuote
				}
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" {
					out = append(out, stmt)
				}
				current.Reset()
				continue
			}
		}
		current.WriteRune(r)
	}

	if tail := strings.TrimSpace(current.String()); tail != "" {
		out = append(out, tail)
	}

	return out
}

func applyStatement(state *schemaState, stmt string) {
	norm := strings.TrimSpace(stmt)
	upper := strings.ToUpper(norm)

	switch {
	case strings.HasPrefix(upper, "CREATE TABLE "):
		applyCreateTable(state, norm)
	case strings.HasPrefix(upper, "ALTER TABLE "):
		applyAlterTable(state, norm)
	default:
		// Ignored: create schema, indexes, extension, comments, etc.
	}
}

func applyCreateTable(state *schemaState, stmt string) {
	re := regexp.MustCompile(`(?is)^CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z0-9_\.]+)\s*\((.*)\)$`)
	m := re.FindStringSubmatch(strings.TrimSpace(stmt))
	if len(m) != 3 {
		return
	}

	tableName := qualifyTableName(m[1])
	body := m[2]
	items := splitTopLevelComma(body)

	t := &table{name: tableName}
	for _, item := range items {
		if col, ok := parseColumnDef(item); ok {
			t.columns = append(t.columns, col)
		}
	}
	state.tables[tableName] = t
}

func applyAlterTable(state *schemaState, stmt string) {
	trimmed := strings.TrimSpace(stmt)

	addRe := regexp.MustCompile(`(?is)^ALTER\s+TABLE\s+([a-zA-Z0-9_\.]+)\s+ADD\s+COLUMN\s+(.+)$`)
	if m := addRe.FindStringSubmatch(trimmed); len(m) == 3 {
		tableName := qualifyTableName(m[1])
		if col, ok := parseColumnDef(m[2]); ok {
			t := ensureTable(state, tableName)
			t.columns = append(t.columns, col)
		}
		return
	}

	dropRe := regexp.MustCompile(`(?is)^ALTER\s+TABLE\s+([a-zA-Z0-9_\.]+)\s+DROP\s+COLUMN\s+([a-zA-Z0-9_]+)$`)
	if m := dropRe.FindStringSubmatch(trimmed); len(m) == 3 {
		tableName := qualifyTableName(m[1])
		colName := strings.ToLower(strings.TrimSpace(m[2]))
		t := ensureTable(state, tableName)
		filtered := make([]column, 0, len(t.columns))
		for _, c := range t.columns {
			if c.name != colName {
				filtered = append(filtered, c)
			}
		}
		t.columns = filtered
		return
	}

	renameTableRe := regexp.MustCompile(`(?is)^ALTER\s+TABLE\s+([a-zA-Z0-9_\.]+)\s+RENAME\s+TO\s+([a-zA-Z0-9_\.]+)$`)
	if m := renameTableRe.FindStringSubmatch(trimmed); len(m) == 3 {
		oldName := qualifyTableName(m[1])
		newName := qualifyTableName(m[2])
		if t, ok := state.tables[oldName]; ok {
			delete(state.tables, oldName)
			t.name = newName
			state.tables[newName] = t
		}
		return
	}

	renameColumnRe := regexp.MustCompile(`(?is)^ALTER\s+TABLE\s+([a-zA-Z0-9_\.]+)\s+RENAME\s+COLUMN\s+([a-zA-Z0-9_]+)\s+TO\s+([a-zA-Z0-9_]+)$`)
	if m := renameColumnRe.FindStringSubmatch(trimmed); len(m) == 4 {
		tableName := qualifyTableName(m[1])
		oldCol := strings.ToLower(strings.TrimSpace(m[2]))
		newCol := strings.ToLower(strings.TrimSpace(m[3]))
		t := ensureTable(state, tableName)
		for i := range t.columns {
			if t.columns[i].name == oldCol {
				t.columns[i].name = newCol
				break
			}
		}
	}
}

func ensureTable(state *schemaState, tableName string) *table {
	if t, ok := state.tables[tableName]; ok {
		return t
	}
	t := &table{name: tableName}
	state.tables[tableName] = t
	return t
}

func splitTopLevelComma(s string) []string {
	var parts []string
	depth := 0
	inSingleQuote := false
	inDoubleQuote := false
	start := 0
	runes := []rune(s)

	for i, r := range runes {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '(':
			if !inSingleQuote && !inDoubleQuote {
				depth++
			}
		case ')':
			if !inSingleQuote && !inDoubleQuote && depth > 0 {
				depth--
			}
		case ',':
			if !inSingleQuote && !inDoubleQuote && depth == 0 {
				parts = append(parts, strings.TrimSpace(string(runes[start:i])))
				start = i + 1
			}
		}
	}

	if start < len(runes) {
		parts = append(parts, strings.TrimSpace(string(runes[start:])))
	}

	return parts
}

func parseColumnDef(def string) (column, bool) {
	trimmed := strings.TrimSpace(def)
	if trimmed == "" {
		return column{}, false
	}

	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "CONSTRAINT ") ||
		strings.HasPrefix(upper, "FOREIGN KEY") ||
		strings.HasPrefix(upper, "PRIMARY KEY") ||
		strings.HasPrefix(upper, "UNIQUE ") ||
		strings.HasPrefix(upper, "CHECK ") {
		return column{}, false
	}

	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		return column{}, false
	}

	colName := normalizeIdent(parts[0])
	rest := strings.TrimSpace(trimmed[len(parts[0]):])
	colType, constraints := splitTypeAndConstraints(rest)
	if colType == "" {
		return column{}, false
	}

	return column{
		name:        colName,
		typ:         colType,
		constraints: constraints,
	}, true
}

func splitTypeAndConstraints(rest string) (string, string) {
	constraintKeywords := []string{
		"NOT", "NULL", "DEFAULT", "PRIMARY", "REFERENCES", "CHECK", "UNIQUE", "CONSTRAINT", "DEFERRABLE", "ON",
	}

	tokens := strings.Fields(rest)
	if len(tokens) == 0 {
		return "", ""
	}

	split := len(tokens)
	for i, tok := range tokens {
		u := strings.ToUpper(tok)
		for _, kw := range constraintKeywords {
			if u == kw {
				split = i
				goto done
			}
		}
	}

done:
	typ := strings.Join(tokens[:split], " ")
	constraints := strings.Join(tokens[split:], " ")
	return strings.TrimSpace(typ), strings.TrimSpace(constraints)
}

func normalizeIdent(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, `"`)
	return strings.ToLower(name)
}

func qualifyTableName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, `"`)
	name = strings.ToLower(name)
	if strings.Contains(name, ".") {
		return name
	}
	return "public." + name
}

func updateDoc(docPath, latest string, state *schemaState, files []string) error {
	content, err := os.ReadFile(docPath)
	if err != nil {
		return err
	}
	doc := string(content)

	doc = updateMigrationHeader(doc, latest)

	for tableName, t := range state.tables {
		heading := "## " + tableName
		if strings.Contains(doc, heading) {
			doc = replaceSectionTable(doc, heading, renderColumnsTable(t.columns))
		}
	}

	doc = replaceMigrationSummary(doc, renderMigrationSummary(files))

	return os.WriteFile(docPath, []byte(doc), 0o644)
}

func updateMigrationHeader(doc, latest string) string {
	re := regexp.MustCompile("as of migration `\\d{4}`")
	return re.ReplaceAllString(doc, fmt.Sprintf("as of migration `%s`", latest))
}

func replaceSectionTable(doc, heading, tableMarkdown string) string {
	start := strings.Index(doc, heading)
	if start == -1 {
		return doc
	}
	sectionStart := start + len(heading)
	next := strings.Index(doc[sectionStart:], "\n## ")
	sectionEnd := len(doc)
	if next != -1 {
		sectionEnd = sectionStart + next
	}

	section := doc[start:sectionEnd]
	re := regexp.MustCompile(`(?ms)^\|.*?\n\|[-| ]+\n(?:\|.*\n)+`)
	if !re.MatchString(section) {
		return doc
	}
	replaced := re.ReplaceAllString(section, tableMarkdown+"\n")

	var b bytes.Buffer
	b.WriteString(doc[:start])
	b.WriteString(replaced)
	b.WriteString(doc[sectionEnd:])
	return b.String()
}

func renderColumnsTable(cols []column) string {
	headers := []string{"Column", "Type", "Constraints"}

	rows := make([][]string, 0, len(cols))
	for _, c := range cols {
		constraints := c.constraints
		if constraints == "" {
			constraints = ""
		}
		rows = append(rows, []string{"`" + c.name + "`", c.typ, constraints})
	}

	return renderMarkdownTable(headers, rows)
}

func replaceMigrationSummary(doc, summary string) string {
	heading := "## Migration History Summary"
	start := strings.Index(doc, heading)
	if start == -1 {
		return doc
	}
	sectionStart := start + len(heading)
	next := strings.Index(doc[sectionStart:], "\n## ")
	sectionEnd := len(doc)
	if next != -1 {
		sectionEnd = sectionStart + next
	}

	prefix := doc[:sectionStart]
	suffix := doc[sectionEnd:]

	return prefix + "\n\n" + summary + "\n" + suffix
}

func renderMigrationSummary(files []string) string {
	headers := []string{"Migration", "Description"}
	rows := make([][]string, 0, len(files))

	for _, file := range files {
		base := filepath.Base(file)
		migration := ""
		description := base
		if len(base) >= 4 {
			migration = base[:4]
		}
		if strings.HasSuffix(base, ".up.sql") && len(base) > 11 {
			raw := base[5 : len(base)-7]
			description = strings.ReplaceAll(raw, "_", " ")
		}
		rows = append(rows, []string{migration, description})
	}

	return renderMarkdownTable(headers, rows)
}

func renderMarkdownTable(headers []string, rows [][]string) string {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i := range headers {
			if i < len(row) && len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	var b strings.Builder
	b.WriteString("|")
	for i, h := range headers {
		b.WriteString(" ")
		b.WriteString(padRight(h, widths[i]))
		b.WriteString(" |")
	}
	b.WriteString("\n")

	b.WriteString("|")
	for i := range headers {
		b.WriteString(strings.Repeat("-", widths[i]+2))
		b.WriteString("|")
	}
	b.WriteString("\n")

	for _, row := range rows {
		b.WriteString("|")
		for i := range headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			b.WriteString(" ")
			b.WriteString(padRight(cell, widths[i]))
			b.WriteString(" |")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
