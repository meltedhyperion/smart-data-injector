## Smart Data Injector

A Data Injection Pipeline tool that uses machine learning (LLM Model) for source to target association mapping and data injection. It utilizes Kestra for orchestrating workflows such as s3 triggers, db injections, ETL Services for data transformation and loading.

### The Idea

![image](https://github.com/user-attachments/assets/4c140a12-1bfa-430b-ac59-a1dcce17b36e)

### One-shot 

<img width="612" alt="image" src="https://github.com/user-attachments/assets/8c79d171-23d1-43a4-995c-46d197c57f4a">


### The workflow diagram

![image](https://github.com/user-attachments/assets/464868ca-b614-4bd5-9b6e-8bade9a60f10)


### Requirements

Golang installed on your system.
- To run the backend,  create `.env` from `.env.example` use command:
  
  ```
  go mod tidy
  go run ./cmd
  ```

- To start the Kestra platform locally, have docker setup in your system, and then use this script:

  ```
  source run_kestra.sh 
  ```
  - On your browser open `localhost:8080/`
  - Create New Flow
  - copy the flow configuration from [here](https://github.com/meltedhyperion/smart-data-injector/blob/main/.kestra_config/my_flow.yaml)
  - add your AWS, MONGODB and GEMINI keys into variables.
  - Save the flow.
  
- To run the client in dev mode, create `.env` in /client directory and use pnpm package manager:

  ```
  cd client
  pnpm install
  pnpm dev
  ```
  - On your browser open `localhost:3000/`

### Execution
1) Go for generate API key.
<img width="781" alt="image" src="https://github.com/user-attachments/assets/d7b8d9c3-6d48-4c24-a4b5-3923af1ecaee">



2) Use the API Key for data injection
<img width="781" alt="image" src="https://github.com/user-attachments/assets/7ff1897f-d58d-4240-92af-a42192e59de8">



3) Once completed, proceed to analytics dashboard:

 <img width="484" alt="image" src="https://github.com/user-attachments/assets/00ddb053-31bb-4d80-b2f2-9ed9d995516c">


## Pipeline task execution log format

<img width="669" alt="image" src="https://github.com/user-attachments/assets/3c3b9b5a-cb14-4c36-a324-5e2738b6055e">
