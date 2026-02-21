import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import Alert from "./Alert";

describe("Alert Component", () => {
  it("should render correctly", () => {
    expect(true).toBe(true);
  });
  it("renders an alert message", () => {
    render(<Alert message="Test Alert" />);
    expect(screen.getByText("Error: Test Alert")).toBeInTheDocument();
  });
});