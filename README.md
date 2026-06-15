## Backend Development Progress

### Architecture

```text
Client
  ↓
API Gateway
  ↓
AWS Lambda (Go)
  ↓
OpenAI API
  ↓
JSON Parsing
  ↓
Supabase PostgreSQL
```

### Implemented Features

#### AWS Serverless API

* Built REST API using Go
* Deployed backend logic to AWS Lambda
* Connected API Gateway to Lambda endpoints
* Configured environment variables and secret management

#### OpenAI Integration

* Integrated OpenAI API using the official Go SDK
* Implemented structured JSON responses using JSON Schema
* Generated TOEIC Part 5 vocabulary quiz data through AI

#### Quiz Data Processing

* Parsed AI-generated quiz data into Go structs
* Validated response format before persistence
* Implemented repository and service layer architecture

#### Supabase Integration

* Connected Lambda to Supabase PostgreSQL
* Designed quiz storage schema
* Implemented batch insertion logic for:

  * `test_quizzes`
  * `test_quiz_questions`
  * `test_quiz_choices`
  * `test_quiz_explanations`

#### Database Design

Quiz data is normalized and stored separately to improve scalability and maintainability.

```text
test_quizzes
    └─ test_quiz_questions
            ├─ test_quiz_choices
            └─ test_quiz_explanations
```

### Current Status

✅ Go API Server Development

✅ AWS Lambda Deployment

✅ API Gateway Integration

✅ OpenAI API Integration

✅ Structured JSON Response Parsing

✅ Supabase PostgreSQL Connection

✅ Quiz / Choice / Explanation Data Persistence

🔄 Improving word selection logic using existing vocabulary data

🔄 Transaction handling and production-level error recovery
