# ml-crawler

This project is a web-crawler responsible to retrieve information about an item from one of the biggest e-commerce platform called MercadoLivre. To developed this web-crawler the following libraries were employed:

```
1. colly - to extract information from the e-commerce platform;
2. strutil - provides string metrics for calculating string similarity as well as other string utility functions.
```


The web-crawler works as one of the agents' actuators, which is responsible to retrieve information about some specific item. Following is some quick explanation of how it can works during the agent's reasoning cycle:
1. Agent's sensors retrieve data from the web about a certain item. 
2. Neural network can learn and improve an existent plan, which could improve an agent's strategy.

For more information and details about the developed agents, one can access our paper here: https://www.researchgate.net/publication/358595795_A_Mediator_Agent_based_on_Multi-Context_System_and_Information_Retrieval

## Run

```
$ go build -o ml
$ ./ml  -term=smartphone
```


## Makefile

There is also a simple makefile that can be used instead of go commands. If you use only make or make all, it is required to pass the search term as a variable. 

```
$ make ITEM=smartphone

```
