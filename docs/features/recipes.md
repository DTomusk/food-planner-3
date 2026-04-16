# Recipes
This document describes how recipes are created, modified, and stored. 

## Recipe model and versioning 
- Recipes are formed of containers and versions 
- Containers store immutable facts about a recipe, e.g. creation time and owner 
- A container can have many versions and each version belongs to one container 
- Versions store the data of recipes that can be edited. This includes name, description, ingredients, image etc. 
- All updates are posts. You cannot edit an existing version, you must create a new one. This means that recipe versions are immutable as well. 

## Recipe ingredients 
- Ingredients in recipes are stored as ingredient usages. These combine an ingredient id with a quantity and a unit. 