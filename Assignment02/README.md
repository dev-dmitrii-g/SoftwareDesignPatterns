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
## Part B. Abstract Factory
### 1. Abstract Products
- `Pokemon` - abstract product class that defines common state (name, type, hp, attack) and combat methods
- `SpecialMove` - abstract product class that defines move execution and name getter
### 2. Concrete Products
- Fire family - `Magby` paired with `Ember`
- Water family - `Slowbro` paired with `Water Pulse`
- Electric family - `Zapdos` paired with `Thunder Shock`
### 3. Abstract Factory 
Abstract factory interface declaring creation methods for product families (`createPokemon()` and `createSpecialMove()`)
### 4. Concrete Factories
`FireFactory`, `WaterFactory`, `ElectricFactory` - concrete factories that implement the abstract factory interface
### 5. Client Orchestration
Client orchestrates the creation of pokemon and special moves using the abstract factory.