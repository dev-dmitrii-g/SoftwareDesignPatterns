# SDP Assignment 2

## Overview
- Theme: Game World Kit
- Domain: Pokemon 

## Part A. Factory Method 
### 1. Product (`Pokemon`)
An abstract product class that defines common state (name, type, hp, attack) and 2 abstract methods
- `speak()` - abstract method that defines what pokemon should speak
- `attack(Pokemon pokemon)` - abstract method that defines how the pokemon attacks 
### 2. Concrete Products
- `Magby` - Fire type pokemon implementing `EMBER`
- `Slowbro` - Water type pokemon implementing `WATER PULSE`
- `Zapdos` - Electric type pokemon implementing `THUNDER SHOCK`
### 3. Factory (`Pokemon Factory`)
Declares the factory method alongside the business method:
- `abstract Pokemon createPokemon()` - the factory method overriden by subclasses
- `public Pokemon spawn()` - high-level business method that instantiates a pokemon using `createPokemon()`, displays its stats and triggers `speak()`
### 4. Concrete Factories
- `MagbyFactory` - creates `Magby` pokemon
- `SlowbroFactory` - creates `Slowbro` pokemon
- `ZapdosFactory` - creates `Zapdos` pokemon
### 5. Client
Can only interact with the `Pokemon Factory` abstractions. It spawns wild pokemon and runs battle interactions.