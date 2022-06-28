# ml-crawler

This project is a web-crawler responsible for retrieving information about an item from one of the biggest e-commerce platforms, MercadoLivre. To develop this web-crawler the following libraries were employed:

```
1. colly - to extract information from the e-commerce platform;
2. strutil - provides string metrics for calculating string similarity as well as other string utility functions.
```


The web-crawler works as one of the agents' actuators, responsible for retrieving information about some specific item. Following is a quick explanation of how it can work during the agent's reasoning cycle:
```
1. Agent's sensors retrieve data about a specific item from the web. 
2. Neural network can learn and improve an existing plan, which could improve an agent's strategy.
```
For more information and details about the developed agents, one can access our paper here: https://www.researchgate.net/publication/358595795_A_Mediator_Agent_based_on_Multi-Context_System_and_Information_Retrieval

## Run

```
$ go build -o ml
$ ./ml  -term=smartphone
```


## Makefile

A simple makefile can be used instead of go commands. If you use only make or make all, it is required to pass the search term as a variable. 

```
$ make ITEM=smartphone

```
