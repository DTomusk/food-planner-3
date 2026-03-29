# Backlog

## Ready 

- [ ] Add images to recipes
    - Priority: high
    - Area: recipes, fullstack
    - Type: feature
    - Why: it's important to know what a dish looks like
    - DoD:
        - recipe version stores image 
        - image appears on listing and detail screens 
        - image infrastructure figured out, provider (at least current) chosen 
        - recipes default to a standard empty image when image not populated 
        - file drop component built and added to recipe form 

## Planning

- [ ] Manual content moderation
    - Priority: high
    - Area: throughout 
    - Type: feature 
    - Why: since anyone can sign up to the site and post free-text (including recipe instructions and usernames, but also more stuff in the future), we should put guardrails in place to prevent offensive and restricted content being shared. This may involve some auto-moderation in the future, but as a first slice, we can add manual content reporting, content being hidden once reports exceed a configured number, and the facilities for devs (and system administrators in the future) to view and action content that's been reported beyond a certain threshold. 
    - DoD: 

- [ ] Private/public recipes 
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Recipe rating 
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Commenting 
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Meal planning (this will likely be several tickets) 
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Add description to recipes
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD:  

- [ ] Dietary tags (gf, vegan etc.)
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD:     

- [ ] Ingredient registry where users can suggest missing ingredients
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Forking recipes
    - Priority: low
    - Area: recipes
    - Type: feature
    - Why: people on recipe sites often comment that they made the recipe with a billion substitutions. Users should be able to save versions of recipes with their own special changes 
    - DoD: 
        - user can copy a specific recipe version 
        - they can make changes 
        - UI shows that it's a fork of a given recipe with a link to the parent 

- [ ] Are you sure you want to exit form
    - Priority: medium
    - Area: recipes, frontend
    - Type: feature
    - Why: users can lose a lot of progress by accidentally clicking off a form 
    - DoD: 
        - Users get a warning whenever they perform an action that would take them off a form page 
            - clicking back 
            - back (mouse)
            - closing tab/window 

## Item template 

- [ ] Item name
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

## Done

- [x] Nav bar item styling 
    - Priority: medium
    - Area: UI, navigation
    - Type: feature
    - Why: currently, the nav bar items suck, they don't look good. We should make them look better
    - DoD: 
        - nav items look better 
        - mobile considered

- [x] Recipe search  
    - Priority: medium
    - Area: recipes fullstack
    - Type: feature
    - Why: you can't currently search for recipes, this is atrocious
    - DoD: 
        - graphql query for searching recipe 
        - ranked search results by relevance 
        - paginated results 
        - recipes/search page 
        - search bar component (maybe on home page)

- [x] Ingredient search 
    - Priority: high 
    - Area: frontend (for now)
    - Type: feature 
    - Why: currently, there are very few ingredients to choose from, so it's not a problem that you can't search them. However, this will VERY quickly stop being the case, and the site is essentially useless without it. When picking ingredients, users should only see ones that are relevant to their search
    - DoD:
        - Ingredient select is searchable, users can type in characters 
        - Search hook introduced: current ingredients are filtered by the given search characters and top five or so matches returned 
        - More ingredients added (doesn't have to be this ticket, but it makes sense to do it here)
        - Nice to have: cache ingredients on the frontend so they don't have to be fetched every time the edit page is opened (react query might already do this for us? Investigate)

- [x] Investigate test coverage tools
    - Priority: high
    - Area: frontend and backend
    - Type: tech debt
    - Why: there is currently no visibility on how much of the code is covered by the current suite of tests, meaning it's much easier to not implement tests and especially test cases 
    - DoD: 
        - Test strategy documented for frontend and backend 
        - Tool found for showing code coverage status (or determined to be too difficult/expensive to do now)
        - AI focused test generation docs created 

- [x] Refresh tokens 
    - Priority: high
    - Area: auth backend (some frontend)
    - Type: feature/tech debt
    - Why: JWTs are short-lived by design, so users cannot stay logged in persistently with just a JWT. We need to introduce sessions so that users can stay logged in on a device for a configured amount of time before having to reauthenticate. 
    - DoD:
        - User can stay logged in longer than their JWT allows 
        - Refresh token stored in http only cookie (so cannot be accessed by js)
        - All tokens in a family invalidated when an expired token is used 
        - Tokens rotated when new ones issued
        - Token hashes stored in db to prevent problems if db leaks

- [x] Log out on unauthorized response 
    - Priority: high
    - Area: auth frontend
    - Type: bug
    - Why: right now, the frontend keeps thinking a user is logged in as long as there's a jwt in local storage. This means that users can access restricted pages while not technically logged in, it's just the API will return unauthorised. We need to remove the jwt from local storage when an unauthorised response comes back.
    - DoD: 
        - JWT removed from local storage when API returns unauthorised 
        - User doesn't see logged in options or protected routes after getting an unauthorised response and before logging in again 

- [x] Markdown formatting
    - Priority: medium
    - Area: frontend
    - Type: bug
    - Why: the input for recipe instructions is a markdown editor, but markdown is not rendered correctly. Titles aren't differentiated and bullet points don't appear. 
    - DoD: 

- [x] Set up configurable theming
    - Priority: high
    - Area: frontend
    - Type: tech debt/feature
    - Why: the site currently uses hardcoded colours for what should be themed components. We should be able to change the primary colour and have it propagate throughout the site, rather than editing magic strings everywhere. We can't have an established design language if we're hardcoding utility classes everywhere without any view of consistency
    - DoD: 
        - Dev can update primary and secondary colour in one place (SSoT)

- [x] Introduce storybook
    - Priority: medium
    - Area: frontend
    - Type: tech debt
    - Why: we now have a lot of components that are frustrating to design and test on the site because you have to find the right pages and have enough data to test them
    - DoD: 
        - Storybook installed 
        - Stories created for components 

- [x] Build out docs area
    - Priority: High
    - Area: ./docs
    - Type: Tech debt
    - Why: Decisions aren't being documented yet, someone checking out the repo won't know what's going on and neither will I a couple months down the line. AI agents now have to read through the repo every time there's a prompt, which slows things down a great deal and increases the chance that they generate code inconsistent with the style and architecture of the repo 
    - DoD: 
        - docs folder created 
        - current scope of repro documented, including architecture and features 
        - AI-optimised area created