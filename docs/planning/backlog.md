# Backlog

## Ready

- [ ] Build out docs area
    - Priority: High
    - Area: ./docs
    - Type: Tech debt
    - Why: Decisions aren't being documented yet, someone checking out the repo won't know what's going on and neither will I a couple months down the line. AI agents now have to read through the repo every time there's a prompt, which slows things down a great deal and increases the chance that they generate code inconsistent with the style and architecture of the repo 
    - DoD: 
        - docs folder created 
        - current scope of repro documented, including architecture and features 
        - AI-optimised area created

## Planning

- [ ] Introduce storybook
    - Priority: medium
    - Area: frontend
    - Type: tech debt
    - Why: we now have a lot of components that are frustrating to design and test on the site because you have to find the right pages and have enough data to test them
    - DoD: 
        - Storybook installed 
        - Stories created for components 

- [ ] Investigate test coverage tools
    - Priority: high
    - Area: frontend and backend
    - Type: tech debt
    - Why: there is currently no visibility on how much of the code is covered by the current suite of tests, meaning it's much easier to not implement tests and especially test cases 
    - DoD: 
        - Test strategy documented for frontend and backend 
        - Tool found for showing code coverage status (or determined to be too difficult/expensive to do now)
        - AI focused test generation docs created

- [ ] Set up configurable theming
    - Priority: high
    - Area: frontend
    - Type: tech debt/feature
    - Why: the site currently uses hardcoded colours for what should be themed components. We should be able to change the primary colour and have it propagate throughout the site, rather than editing magic strings everywhere. We can't have an established design language if we're hardcoding utility classes everywhere without any view of consistency
    - DoD: 
        - Dev can update primary and secondary colour in one place (SSoT)

## Item template 

- [ ] Item name
    - Priority: 
    - Area: 
    - Type: (bug, tech debt, feature)
    - Why: 
    - DoD: 