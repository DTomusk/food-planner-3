# Ingredient taxonomy 
I want to design a system for organising ingredients that could lead to powerful features in the future and that, in the meantime, would make organising and expanding the list of ingredients easier. 

I want to categorise ingredients. The first level of categorisation I have in mind is: 
- Raw ingredients
- Derivatives 
- Composites 

## Raw ingredients
Raw ingredients will be ingredients in their natural form as they are cultivated. This can be a chicken, an apple, coffee beans, etc. The main point being that nothing has been done to these ingredients by people. 

Examples of raw ingredients: 
- Milk (although pasteurised etc. it essentially is how it comes out)
- Chicken (not its constituent components e.g. thigh, breast, bones)
- Nuts (not ground, not almond flour etc.)
- Whole fruits 
- Salt
- Rice

## Derivatives 
These are ingredients that are (mainly) derived from one ingredient by some kind of process. All derivatives have one parent raw ingredient. 

Examples: 
- Flour (milled from wheat)
- Cuts of meat and mince (butchered/ground from animal)
- Egg whites and yolks (simply extracted by hand)
- Oils (pressed or however they're done, but perhaps not falvoured oils I don't know)
- Sugars 
- Butter and cheese (although many cheeses use multiple ingredients and varying processes, they're essentially derived from milk)

## Composites 
These are many ingredients that have been combined in a certain way. 

Examples:
- Spice paste/mix 
- Soy sauce 
- Pasta 
- Bread 
- Noodles
- Sauces 

Note: there won't always be a clear categorisation for an ingredient, so it's not incredibly important for everything to be perfectly in its right place 

This can simply be an organisational categorisation to start with, and later become a graph. 


# Modelling relationships 
Ingredients are related to each other in various ways. 

## Taxonomy relationship
Ingredients are related in terms of specificity. All brown onion are onions, not all onions are brown onions. All Tagliatelle are pasta, not all pasta are tagiatelle. We can call this the taxonomy relationship. Onion is the taxonomy parent of brown onion

## Derivative relationship 
Derivatives are related to the ingredients they're derived from. Flour is milled wheat. An egg yolk is taken from an egg. Beef mince is minced beef. There is a process that transforms a raw ingredient into a derived ingredient. 

## Composite relationship 
Composites contain components. Fresh pasta contains eggs and flour. The component relationship doesn't have to be exhaustive. We should define a limited list of components that communicate a kind of "profile" of what the composite is. Tomato sauce contains tomatoes, could contain a lot of other things. 

# Schema 
I need to start thinking about how to represent this in the yaml reference file and the database schema. 

Here's the current ingredients table schema: 
| Column           | Type    | Constraints          |
|------------------|---------|----------------------|
| `id`             | UUID    | PRIMARY KEY          |
| `name`           | TEXT    | NOT NULL             |
| `preferred_unit` | INTEGER | NOT NULL             |
| `file_key`       | TEXT    | NOT NULL, UNIQUE     |
| `counter`        | TEXT    | NULLABLE             |
| `plural`         | TEXT    | NULLABLE             |
| `counter_plural` | TEXT    | NULLABLE             |

I don't think any existing columns have to be removed or modified. We can add new columns for: 
- Processing level (raw = 1, derivative = 2, composite = 3)
- Taxonomy parent (nullable uuid derived from key string in yaml file)