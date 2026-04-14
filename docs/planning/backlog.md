# Backlog

## Ready

- [ ] Event pub/sub
    - Priority: high
    - Area: backend, infrastructure
    - Type: tech debt
    - Why: in the last story we added audit logging. Currently this happens in the same call as the action itself. However, this is a side effect that can be handled asynchronously, and if we want to add more side effects (which we will, e.g. notifications coming up), current services will be bloated by dependencies, responsibilities will be difficult to understand (e.g. what's in the critical path and needs to happen to return to the client and what's not) and things will be much harder to test and extend. We want to add a thin slice in the direction of having backgrounded workers and event infrastructure. In this first iteration, this can still run on the API service and in memory, we don't need message queues etc. yet, but we want to build it in a way that's easy to extend in the future (which we will need for retries, guaranteed delivery, and horizontal scaling)
    - DoD: 
        - In memory event bus (pub) implemented that existing services can call to emit events
        - Event consumers/handlers implemented that asynchronously handle side effects
        - Audit logging moved to event handler, so the audit log from the previous story (user signed up successfully) gets created outside of the main flow

## Planning 

- [ ] Manual content moderation
    - Depends on: Audit logging
    - Priority: high
    - Area: throughout 
    - Type: feature 
    - Why: since anyone can sign up to the site and post free-text (including recipe instructions and usernames, but also more stuff in the future), we should put guardrails in place to prevent offensive and restricted content being shared. This may involve some auto-moderation in the future, but as a first slice, we can add manual content reporting, content being hidden once reports exceed a configured number, and the facilities for devs (and system administrators in the future) to view and action content that's been reported beyond a certain threshold. 
    - DoD:
        - Users can report content 
            - Limit to recipe versions 
        - Users can select one of a number of reasons (reasons persisted as enum)
        - Users can write in details of their complaint
            - If reason is other than "Other", details are optional 
            - If reason is "Other", then details are required
        - Users can't report the same recipe version more than once
        - System saves each report including:
            - Reporting user
            - Recipe version id and recipe id
            - Reason
            - Optional details
            - Report timestamp 
        - Recipe versions store report count
        - Recipe versions above a configurable report threshold are hidden from non-moderator users
            - Visible to the author and to moderators
        - Moderation events (dismiss, restore, remove, note) are recorded to audit log with actor, action, reason, and timestamp
        - Author is notified that their recipe is hidden
        - Previous versions remain hidden when a new version is created and approved
        - Dev can run scripts or use backend tools to review and resolve reports

- [ ] Manual content moderation
    - Priority: high
    - Area: throughout 
    - Type: feature 
    - Why: since anyone can sign up to the site and post free-text (including recipe instructions and usernames, but also more stuff in the future), we should put guardrails in place to prevent offensive and restricted content being shared. This may involve some auto-moderation in the future, but as a first slice, we can add manual content reporting, content being hidden once reports exceed a configured number, and the facilities for devs (and system administrators in the future) to view and action content that's been reported beyond a certain threshold. 
    - DoD:
        - Users can report content 
            - Limit to recipe versions 
        - Users can select one of a number of reasons (reasons persisted as enum)
        - Users can write in details of their complaint
            - If any reason other than other, details are optional 
            - If Other, then must add details
        - Users can't report the same content more than once
        - System saves each report including:
            - Reporting user
            - Entity id
            - Entity type 
            - Content at time of reporting 
            - Report timestamp 
        - Entities store number of reports 
        - Entities with a number of reports greater than a configured value are hidden from other users 
            - What if a user has saved someone else's recipe?
        - The user who created the entity is informed that it's been hidden?
        - 

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

- [ ] Profile images
    - Priority: medium 
    - Area: users, fullstack
    - Type: feature 
    - Why: users can make their profiles more their own if they can associate a profile image with themselves 
    - DoD:
        - [ ] users can upload profile images using the existing upload infrastructure 
        - [ ] users can change propfile picture 

- [ ] Allergens
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Dish storage time (freezable etc.)
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Code coverage improvements 
    - Priority: 
    - Area: 
    - Type: tech debt
    - Why: right now, code coverage is checking files we can't test (e.g. gql generated files). This brings down the score and gives a false understanding of the actual coverage, making the coverage report much less useful
    - DoD: 

- [ ] Save others' recipes
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] Unit conversion
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 

- [ ] User preferences
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: users can set dietary requirements that apply to all searches by default (but can be changed). Users can choose whether they prefer metric or imperial
    - DoD: 

- [ ] RBAC
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

- [x] Audit logging
    - Priority: high
    - Area: visibility for what's going on in the system
    - Type: tech debt
    - Why: we want to have an audit trail of what's going on in the site. This will record actions carried out. This will be useful for content moderation but also for other things 
    - DoD: 
        - audit_log table implemented with columns: actor_id, action, resource_type, resource_id, old_state, new_state, reason, result, context, created_at
        - audit logging service created that can be called from services to record events
        - user signup audit event implemented and recorded for all new signups
        - dev can query audit log to see signup history

- [x] Dietary tags (gf, vegan etc.)
    - Priority: high
    - Area: recipes, ingredients, fullstack
    - Type: feature
    - Why: it should be easy for people with different dietary requirements to search recipes that meet their needs/wants
    - DoD: 
        - [x] animal product enum on ingredients (i.e. vegan, vegetarian, meat)
        - [x] gluten on ingredients 
        - [x] recipe versions store whether they are whatever 
        - [x] users can filter on dietary tags 
        - [x] recipes show clearly in listing and detail what they are

- [x] Add description to recipes
    - Priority: medium
    - Area: recipes, fullstack
    - Type: feature
    - Why: it would be good for recipes to have a tagline/short description that would appear both in listing and detail views. 
    - DoD: 
        - [x] description added to recipe version 
        - [x] editable in create and update forms 
        - [x] appears in listing card
        - [x] appears in detail

- [x] Are you sure you want to exit form
    - Priority: medium
    - Area: recipes, frontend
    - Type: feature
    - Why: users can lose a lot of progress by accidentally clicking off a form 
    - DoD: 
        - Users get a warning whenever they perform an action that would take them off a form page 
            - [x] clicking back 
            - [x] back (mouse)
            - [x] closing tab/window 

- [x] Proper breadcrumbs
    - Priority: high
    - Area: frontend navigation
    - Type: tech debt
    - Why: navgiation is really bad at the moment. Some pages have a back button that isn't clear about what it does, and some don't have one at all. You can go back to a form for example by clicking back. The back button isn't as useful as breadcrumbs would be
    - DoD: 
        - [x] breadcrumb component implemented
        - [x] hierarchy decided 
        - [x] back button removed

- [x] Add images to recipes
    - Priority: high
    - Area: recipes, fullstack
    - Type: feature
    - Why: it's important to know what a dish looks like
    - DoD:
        - [x] recipe version stores image 
        - [x] image appears on listing and detail screens 
        - [x] image infrastructure figured out, provider (at least current) chosen 
        - [x] recipes default to a standard empty image when image not populated 
        - [x] file drop component built and added to recipe form
        - [x] a process to remove unused images 
        - [x] a process to claim images for entities

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