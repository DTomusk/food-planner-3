# Ingredients 
This document describes how ingredients are stored and updated 

## Persistence 
Ingredients are stored in the database in a reference table. The reference table is populated with a reference file. Beyond the sync process that copies ingredients form the reference file to the database, there is no service etc. that can mutate ingredients. 

## Sync process 
There's a batched sync process that copies ingredient data from the reference file to the ingredient repo. It is a command with an associated dockerfile. The sync process reads the file and extracts ingredients in a file specific structure, converts these to ingredient entities, and upserts them into the ingredient repo. 