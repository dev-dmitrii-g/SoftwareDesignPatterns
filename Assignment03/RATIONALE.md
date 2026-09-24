# Rationale for this project
## Problem I chose
I chose Neural Implants implementation. The main reason is that I saw my friend playing Cyberpunk 2077 and I wanted to implement it.
## Why Bridge or Adapter pattern alone would not be enough
In the scope of this project I used both patterns because we have 2 compatible implementations from major companies and one implementation from a black market.
Since the black market implementation is pretty janky, I needed to wrap it inside the adapter to properly translate it to the target interface. 
The bridge pattern is used to decouple the abstraction from the implementation to prevent the need to implement each permutation of the interface.
## Why wrapped implementation is genuinely incompatible
The black market implementation is not compatible with the target interface because it has a completely different structure under the hood. It accepts a voltage and a byte array of data, instead of the region, intensity and the command. Then it uses a completly different error mechanism, error codes, while the target interface uses exceptions.
## One limitiation of my final design
The black market implementation provides many different error codes, but the target implementation only uses one exception.