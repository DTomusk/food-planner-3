import { renderHook } from '@testing-library/react'
import { useFormatIngredients } from './useFormatIngredients'
import { describe, expect, test } from 'vitest'

describe("formatIngredient", () => {
  const { result } = renderHook(() => useFormatIngredients())

  test.each([
    [
      "formats ingredient with unit symbol",
      {
        name: "Test Ingredient",
        quantity: 2,
        unitSymbol: "g",
        counter: null,
        plural: null,
        counterPlural: null,
      },
      "2g test ingredient",
    ],
    [
      "formats ingredient without unit",
      {
        name: "Egg",
        quantity: 3,
        unitSymbol: "",
        counter: null,
        plural: "Eggs",
        counterPlural: null,
      },
      "3 eggs",
    ],
    [
      "formats singular correctly",
      {
        name: "Egg",
        quantity: 1,
        unitSymbol: "",
        counter: null,
        plural: "Eggs",
        counterPlural: null,
      },
      "1 egg",
    ],
    [
        "formats ingredient with counter",
        {
            name: "Bread",
            quantity: 2,
            unitSymbol: "",
            counter: "slice",
            plural: null,
            counterPlural: "slices",
        },
        "2 slices of bread",
    ],
    [
        "formats singular ingredient with counter correctly",
        {
            name: "Bread",
            quantity: 1,
            unitSymbol: "",
            counter: "loaf",
            plural: null,
            counterPlural: "loaves",
        },
        "1 loaf of bread",
    ]
  ])("%s", (_, input, expected) => {
    expect(result.current.formatIngredient(input)).toBe(expected)
  })
})