# Dive Game Documentation
## How to run
- Build the docker file


    docker build -t dive-games-app .
- Run the docker container


     docker run -p 8080:8080 dive-games-app

- Call the service from a browser calling this url (postman or curl can be used as well)


     http://localhost:8080/api/v1/ltp?pairs=BTC/USD&pairs=BTC/EUR

## New Notes
After receiving feedback, I explored different ways to structure the project and decided to use the Standard Layout as defined [here](https://github.com/golang-standards/project-layout). I removed domain-specific errors and reduced dependencies.


It's specified that the price should be numeric without specifying the required number of decimals. Given the sensitivity of this information (even slight changes in cents can significantly impact prices), I've decided to return the exact value as provided by Kraken.

## Requirements
### Explicit Requirements
Create a GET endpoint at '/api/v1/ltp' that will receive one or more pairs and will return the last traded price for them
```json
{
  "ltp": [
    {
      "pair": "BTC/CHF",
      "amount": "49000.12"
    },
    {
      "pair": "BTC/EUR",
      "amount": "50000.12"
    },
    {
      "pair": "BTC/USD",
      "amount": "52000.12"
    }
  ]
}
```
### Missing Requirements
- What action should be taken if some pairs are found while others are not? What should be return if none is found?
    
    
    In this case I decided to fail if any or all of the values cannot be found or if any error is found in the runtime
- How many requests per minute (rpm) are expected? Can we (functional wise) cache the values for 60 seconds? Would we need it?
    

    I'm assuming it's the first iteration of the delivery, therefore I'm leaving such considerations for the future

## Design Choices
- I implemented the Clean Architecture strategy to define the packages


- I chose to use testify dependency for clarity in the test assertions


- I've implemented standard HTTP retry mechanisms, timeouts, management of idle connections, and more to enhance resilience 


- I decided to use the PairValues array with the values requested, instead of using the values directly. By doing this, we can add any PairValue without modifying the existing logic. In the future we can decide to read this values from other sources, such as a properties file, a bucket, or a database in the future. By doing so, we adhere to the SOLID principle of being open for extension and closed for modification.


- In general, the SOLID design principles were implemented in the project, including the Single Responsibility Principle and the Dependency Inversion Principle. 
## Testing
- End-to-end integration tests can be found in the main package


- Unit tests reside within their corresponding packages. Only selected tests were implemented, focusing on those that provided significant value. In a real-world scenario, test coverage should ideally exceed 80%

## Considerations
### Interfaces
I created an interface for the KrakenService to remain agnostic of the solution or provider of the information. I did not find it necessary to create any other interfaces
### DTOs
I used DTOs to separate domain information from external resources and responses

The reason behind this is external resources might change in the future, using DTOs to handle their responses make our services and handlers agnostic from them

Using DTOs in the responses allow us to handle the full extent of the trades. Although the response may require less information than what is handled, having the full extent of the data might be useful for other services or endpoints, and therefore we should use them as a whole

The use of DTOs respects SRP principle 
### Naming
I avoided adding redundant information in names. For instance, appending 'dto' to DTOs would unnecessarily lengthen them, as their package name already includes 'dto'
### Errors
I created domain errors so that we can handle each error as desired
### Git
With only one implemented use case, the repository contains just a single commit. I prefer having one commit for each new feature or fix. This practice ensures a clear project history and simplifies the process of rollback
## TODOs
- If the application grows beyond its current state, we should consider implementing a Circuit Breaker to prevent overloading the API


- Kraken APIs also have a set rate limit; we should implement the same mechanism


- Depending on the application's use, we might consider utilizing tools like OpenAPI to define our resources and clients. These tools would simplify our consumers' tasks by allowing them to autogenerate entities and connectors from a YAML file.
