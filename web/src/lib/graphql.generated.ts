import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type AuthPayload = {
  __typename?: 'AuthPayload';
  jwt: Scalars['String']['output'];
  user: User;
};

export type CreateIngredientUsageInput = {
  ingredientID: Scalars['ID']['input'];
  quantity: Scalars['Float']['input'];
  unit: Scalars['Int']['input'];
};

export type CreateRecipeInput = {
  cookMins: Scalars['Int']['input'];
  ingredientUsages: Array<CreateIngredientUsageInput>;
  name: Scalars['String']['input'];
  portions: Scalars['Int']['input'];
  prepMins: Scalars['Int']['input'];
};

export type Ingredient = {
  __typename?: 'Ingredient';
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  preferredUnit: Unit;
};

export type IngredientUsage = {
  __typename?: 'IngredientUsage';
  id: Scalars['ID']['output'];
  ingredient: Ingredient;
  quantity: Scalars['Float']['output'];
  unit: Unit;
};

export type Mutation = {
  __typename?: 'Mutation';
  _empty?: Maybe<Scalars['String']['output']>;
  createRecipe: Recipe;
  signin: AuthPayload;
  signup: AuthPayload;
};


export type MutationCreateRecipeArgs = {
  input: CreateRecipeInput;
};


export type MutationSigninArgs = {
  input: SignInInput;
};


export type MutationSignupArgs = {
  input: SignUpInput;
};

export type Query = {
  __typename?: 'Query';
  _empty?: Maybe<Scalars['String']['output']>;
  ingredients: Array<Ingredient>;
  recipe?: Maybe<Recipe>;
  recipes: Array<Recipe>;
};


export type QueryRecipeArgs = {
  id: Scalars['ID']['input'];
};

export type Recipe = {
  __typename?: 'Recipe';
  cookMins: Scalars['Int']['output'];
  id: Scalars['ID']['output'];
  ingredientUsages: Array<IngredientUsage>;
  name: Scalars['String']['output'];
  portions: Scalars['Int']['output'];
  prepMins: Scalars['Int']['output'];
  user: User;
};

export type SignInInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type SignUpInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Unit = {
  __typename?: 'Unit';
  name: Scalars['String']['output'];
  symbol: Scalars['String']['output'];
  val: Scalars['Int']['output'];
};

export type User = {
  __typename?: 'User';
  email: Scalars['String']['output'];
  id: Scalars['ID']['output'];
};

export type SignInMutationVariables = Exact<{
  input: SignInInput;
}>;


export type SignInMutation = { __typename?: 'Mutation', signin: { __typename?: 'AuthPayload', jwt: string } };

export type SignUpMutationVariables = Exact<{
  input: SignUpInput;
}>;


export type SignUpMutation = { __typename?: 'Mutation', signup: { __typename?: 'AuthPayload', jwt: string } };

export type GetIngredientsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetIngredientsQuery = { __typename?: 'Query', ingredients: Array<{ __typename?: 'Ingredient', name: string, id: string, preferredUnit: { __typename?: 'Unit', name: string, val: number, symbol: string } }> };

export type CreateRecipeMutationVariables = Exact<{
  input: CreateRecipeInput;
}>;


export type CreateRecipeMutation = { __typename?: 'Mutation', createRecipe: { __typename?: 'Recipe', id: string, name: string } };

export type GetRecipeQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetRecipeQuery = { __typename?: 'Query', recipe?: { __typename?: 'Recipe', name: string, id: string, prepMins: number, cookMins: number, portions: number, user: { __typename?: 'User', id: string }, ingredientUsages: Array<{ __typename?: 'IngredientUsage', quantity: number, unit: { __typename?: 'Unit', name: string }, ingredient: { __typename?: 'Ingredient', name: string } }> } | null };

export type GetRecipesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetRecipesQuery = { __typename?: 'Query', recipes: Array<{ __typename?: 'Recipe', name: string, id: string }> };


export const SignInDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SignIn"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"SignInInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"signin"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"jwt"}}]}}]}}]} as unknown as DocumentNode<SignInMutation, SignInMutationVariables>;
export const SignUpDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SignUp"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"SignUpInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"signup"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"jwt"}}]}}]}}]} as unknown as DocumentNode<SignUpMutation, SignUpMutationVariables>;
export const GetIngredientsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetIngredients"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ingredients"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"preferredUnit"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"val"}},{"kind":"Field","name":{"kind":"Name","value":"symbol"}}]}}]}}]}}]} as unknown as DocumentNode<GetIngredientsQuery, GetIngredientsQueryVariables>;
export const CreateRecipeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateRecipe"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateRecipeInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createRecipe"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]} as unknown as DocumentNode<CreateRecipeMutation, CreateRecipeMutationVariables>;
export const GetRecipeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetRecipe"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"recipe"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"ingredientUsages"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"quantity"}},{"kind":"Field","name":{"kind":"Name","value":"unit"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"ingredient"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"prepMins"}},{"kind":"Field","name":{"kind":"Name","value":"cookMins"}},{"kind":"Field","name":{"kind":"Name","value":"portions"}}]}}]}}]} as unknown as DocumentNode<GetRecipeQuery, GetRecipeQueryVariables>;
export const GetRecipesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetRecipes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"recipes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<GetRecipesQuery, GetRecipesQueryVariables>;