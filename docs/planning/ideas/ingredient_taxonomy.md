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

### Note 
I've since rethought this a bit. I think that things that are minimally process (e.g. a segment of a satsuma, an egg yolk, a cut of meat) should still count as raw. It's more significant processing such as mincing, pressing etc. that changes the fundamental qualities of the ingredient that make it derived. 

## Composite relationship 
Composites contain components. Fresh pasta contains eggs and flour. The component relationship doesn't have to be exhaustive. We should define a limited list of components that communicate a kind of "profile" of what the composite is. Tomato sauce contains tomatoes, could contain a lot of other things. 

# Schema 
I need to start thinking about how to represent this in the yaml reference file and the database schema. 

Here's the current ingredients table schema: 
| Column                 | Type    | Constraints      |
|------------------------|---------|------------------|
| `id`                   | UUID    | PRIMARY KEY      |
| `name`                 | TEXT    | NOT NULL         |
| `preferred_unit`       | INTEGER | NOT NULL         |
| `file_key`             | TEXT    | NOT NULL, UNIQUE |
| `counter`              | TEXT    | NULLABLE         |
| `plural`               | TEXT    | NULLABLE         |
| `counter_plural`       | TEXT    | NULLABLE         |
| `animal_product_level` | INTEGER | NOT NULL, DEFAULT 0 |
| `contains_gluten`      | BOOLEAN | NOT NULL, DEFAULT FALSE |

I don't think any existing columns have to be removed or modified. We can add new columns for: 
- Processing level (raw = 1, derivative = 2, composite = 3)
- Taxonomy parent (nullable uuid derived from key string in yaml file)
- Show in search (non nullable flag default to true to show ingredients that can be selected in recipes)

We can also consider adding new tables for derivative parents and components, but that can come later. 

# Supporting multiple units
Currently, we store a preferred unit. The idea behind it was that all other units for the ingredient could resolve down to the preferred unit, which would be used for shopping list purposes, but that doesn't matter until we have shopping lists. In the context of recipes, any acceptable unit is fine. We're not going to be combining recipes to form shopping lists any time soon, but when we do, we're going to need to have a mechanism for converting to a standard "shopping unit" (which I suppose is what the preferred unit represents). 

The question is, do we explicitly store each unit for an ingredient, or do we have a more flexible approach? Storing each possible unit for an ingredient seems stupid. For example, there's no point storing grams as well as kilograms, pounds and ounces because they're all mass measurements. Similarly, teaspoons, millilitres, cups, pints etc. are all easy to convert between. Conversion between mass and volume requires density, but do we need conversion yet? We could allow grams of flour and cups of flour at the moment without any need to convert between them, because recipes exist on their own. If we had a metric/imperial toggle, then we would have to be able to do certain conversions, but in that case we would also have to consider mass/volume because flour is often cups in imperial and grams in metric, so that's a two step process. 

Ok, so the question really is, what units do we allow for an ingredient? A preferred unit can just be the unit that's selected by default, but we could list a lot more. Maybe at the top of the ingredients section when building a recipe we have a metric/imperial toggle which narrows down what units get shown? There's no a priori reason why we couldn't mix unit types in the context of one recipe. We could have 10g butter and 2 cups flour. Maybe units need to have a richer profile. For useability, we shouldn't flood the unit selector with tons of units, we should suggest a sensible default (especially taking a user's preferences into account), and we should allow them to choose what they want.

The next step might be to consider showing more units. For the quantum unit (still can't think of a good name), we may still choose to not show any units, I think that would be simpler for now. You could have half an onion but not 50g onion. And I think for these kinds of ingredients, it makes sense. The preferred unit choice has to be intuitive and there is a cost to changing it until we allow multiple units. Recipes might get into a weird state if the preferred unit changes under them. Now, that shouldn't be the case in a well-built system. An ingredient usage in a recipe can have any ingredient it wants and that should remain valid regardless of what happens to the ingredient it's referencing. 

# Monitoring and updating 
I obviously won't get all the ingredients there are in the first pass. The problem is, it's a very manual process at the moment, and I'm bound to get things wrong. 

# Tasks
- Showing multiple units in the recipe form
    - Slice one might be to have the preferred unit selected and have the unit dropdown show all units of the same type (dimension)
    - Implement dimension metadata on units e.g. count, mass, volume, length
    - Infer what units can be shown. For now, if unit dimension is count, then don't show any units. Else, show preferred unit at the top of list and then all other 
