# Backlog

## Ready

- [ ] Investigate test coverage tools
    - Priority: high
    - Area: frontend and backend
    - Type: tech debt
    - Why: there is currently no visibility on how much of the code is covered by the current suite of tests, meaning it's much easier to not implement tests and especially test cases 
    - DoD: 
        - Test strategy documented for frontend and backend 
        - Tool found for showing code coverage status (or determined to be too difficult/expensive to do now)
        - AI focused test generation docs created

## Planning

- [ ] Introduce storybook
    - Priority: medium
    - Area: frontend
    - Type: tech debt
    - Why: we now have a lot of components that are frustrating to design and test on the site because you have to find the right pages and have enough data to test them
    - DoD: 
        - Storybook installed 
        - Stories created for components 

- [ ] Set up configurable theming
    - Priority: high
    - Area: frontend
    - Type: tech debt/feature
    - Why: the site currently uses hardcoded colours for what should be themed components. We should be able to change the primary colour and have it propagate throughout the site, rather than editing magic strings everywhere. We can't have an established design language if we're hardcoding utility classes everywhere without any view of consistency
    - DoD: 
        - Dev can update primary and secondary colour in one place (SSoT)

- [ ] Ingredient search 
    - Priority: high 
    - Area: frontend (for now)
    - Type: feature 
    - Why: currently, there are very few ingredients to choose from, so it's not a problem that you can't search them. However, this will VERY quickly stop being the case, and the site is essentially useless without it. When picking ingredients, users should only see ones that are relevant to their search
    - DoD:
        - Ingredient select is searchable, users can type in characters 
        - Search hook introduced: current ingredients are filtered by the given search characters and top five or so matches returned 
        - More ingredients added (doesn't have to be this ticket, but it makes sense to do it here)
        - Nice to have: cache ingredients on the frontend so they don't have to be fetched every time the edit page is opened (react query might already do this for us? Investigate)

- [ ] Log out on unauthorized response 
    - Priority: high
    - Area: auth frontend
    - Type: bug
    - Why: right now, the frontend keeps thinking a user is logged in as long as there's a jwt in local storage. This means that users can access restricted pages while not technically logged in, it's just the API will return unauthorised. We need to remove the jwt from local storage when an unauthorised response comes back.
    - DoD: 
        - JWT removed from local storage when API returns unauthorised 
        - User doesn't see logged in options or protected routes after getting an unauthorised response and before logging in again 

- [ ] Refresh tokens 
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

- [ ] Recipe search  
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

## Item template 

- [ ] Item name
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

## Done 

- [x] Build out docs area
    - Priority: High
    - Area: ./docs
    - Type: Tech debt
    - Why: Decisions aren't being documented yet, someone checking out the repo won't know what's going on and neither will I a couple months down the line. AI agents now have to read through the repo every time there's a prompt, which slows things down a great deal and increases the chance that they generate code inconsistent with the style and architecture of the repo 
    - DoD: 
        - docs folder created 
        - current scope of repro documented, including architecture and features 
        - AI-optimised area created