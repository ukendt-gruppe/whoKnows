# How are you DevOps?

We have tried to work with as many DevOps principles as possible throughout the semester. It has been both rewarding and challenging - skewing a bit towards the latter

## Collaboration

### **The "Three Ways" principles**
We are inspired by the Three Ways principles:
- **Flow**: Improve development and deployment processes through automation and clear task allocation
- **Feedback**: Leverage monitoring and metrics to gather actionable insights for system optimization
- **Continuous Learning**: Share knowledge within the team to reduce silos and encourage improvement

### **Version Control**
- Pull Request based approach with feature branches on Github. We tried to implement a development branch for staging which is still a work in progress. This branching strategy allow us to ensure code quality through peer reviews

- When we follow the branching strategy rigorously it encourage collaboration by involving team members in code review discussions. Definitely works better in person

### **Shared Responsibilities & Task Management**
- Collective code-ownership. Worked well in the beginning of the project and a little less well towards the end

- Jira style board in a shared Notion workspace that also contained notes for the individual weeks

- Too ambitious with database setup and manual filtering on tasks for specific weeks made this a not-optimal way of collaborating

### **Pair Programming**
- When possible, we've tried to collaborate on tasks like CI/CD pipelines, db migrations etc. to reduce knowledge silos and encourage knowledge sharing

## Automation

### **CI/CD Pipelines**
- Automated testing, building, package management and deployment using GitHub Actions
- Collective secret management
- Unit tests and integration tests ensure no changes are merged without passing tests

### **Database Management**
- Automated migrations, backups, and restores, minimizing manual errors
- Infrastructure for database deployment managed with Terraform Modules

### **Monitoring and Metrics Collection**
- Making decisions based on the best available data
- Scheduled Postman Monitor for SRE with testing
- Initialy we monitored the containers running on the server with scripts. This was not optimal
- Implemented Prometheus to collect data on system performance
- Visualized through Grafana dashboards

### **Docker**
- Containerized applications ensure consistent environments for developers and in production

## Code Quality

### **Security Findings**
- Security flaws were given top priority and were mostly addressed, although it was sometimes forgotten among other concerns

- Critical vulnerabilities were fixed to prevent potential system risks, although as we had the old code repository in our repository, the improved Python code also muddied the picture, as this was not on our radar, but we forgot to add it to an ignore list

### **Code Repetition and Maintainability**
- Code repetition tolerance was set at a higher percentage for new code

- Maintainability issues were deprioritized due to more urgent requirements

- We acknowledge that deferring maintainability might cause future technical debt

### **Implementation Insights**
- Code quality tools were integrated retrospectively
- More effort was put into fixing issues in newly written code
- Vice versa old code was a bit forgotten when it came to code quality
- An earlier implementation would have been beneficial, but now we're up and running

### **Future Recommendations**
- Implement code quality tools earlier in development
- A culture of immediately addressing the findings 
- Balance immediate requirements with long-term code health
- Regularly review and address technical debt
- Gradually improve code quality across existing and new codebase

## Monitoring Realization

We've got Prometheus running on the server and Grafana locally for now. 

A very clear message from the monitoring has been that we need to redeploy the Prometheus Service to its own server. Now the metrics are reset whenever we're redeploying the application.

This has also led to conversations about deployment strategies, since it very clearly illustrates that whenever we redeploy our service, we experience downtime.

Based on the monitoring of search, we've identified that we need to update the wiki_scraper setup to scrape articles from an area our users are more interested in. 

## Core Issues

### **Knowledge Silos**
- Limited task rotation in complex areas
- Steep learning curves creating hesitation for broader involvement
- Concentrated expertise leading to unofficial "code ownership"
- Lack of time together as a team and context-switching

### **Mitigations**
- Pair programming has been implemented to encourage collaboration on complex CI/CD tasks, but this requires further reinforcement
- Team members are encouraged to document processes and share knowledge with the team, though this practice can be expanded
- Deep-dive sessions on pipeline components will help reduce the steep learning curve for broader involvement
- Assigning developers to research and explain their work has mitigated this to some extent, but this can be bolstered through mentorship and systematic knowledge transfer
- Fostering a high-trust culture where experimentation and mistakes are seen as learning opportunities will encourage a more collaborative approach
