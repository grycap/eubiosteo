PERGAMO
===================

Project for the H2020 Indigo-DataCloud project. Serverless job scheduler on Golang to run docker jobs on demand. The service is accessible through an API-Rest and a frontend webpage in React-Redux.

It has two deployment versions, locally or distributed. When deployed locally, it uses local docker daemon to run the containers. When used in a distributed way, it relies on a slurm cluster deployed with Infrastructure-Manager (IM).
