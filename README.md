# Appointment Service

Service responsible to manage the appointments

# Local Development

## Requirements

- [Kubernetes](https://kubernetes.io/)
- [AWS CLI](https://aws.amazon.com/cli/)

## Manual deployment

### Attention

Before deploying the service, make sure to set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.

Be aware that this process will take a few minutes (~4 minutes) to be completed.

To deploy the service manually, run the following commands in order:

```bash
make init
make check # this will execute fmt, validate and plan
make apply
```

To destroy the service, run the following command:

```bash
make destroy
```

## Automated deployment

The automated deployment is triggered by a GitHub Action.

# Endpoints

Legend:
- ✅: Development completed
- 🚧: In progress
- 💤: Not started


| Completed | Method | Endpoint                                              | Description                       | User Role      |
| --------- | ------ | ----------------------------------------------------- | --------------------------------- | -------------- |
| ✅         | POST   | `/appointments`                                       | Create an appointment             | Patient        |
| ✅         | GET    | `/appointments`                                       | Get all appointments              | Doctor/Patient |
| ✅         | GET    | `/appointments/{appointmentId}`                       | Get an appointment by id          | Doctor/Patient |
| ✅         | PUT    | `/appointments/{appointmentId}`                       | Update an appointment             | Patient        |
| ✅         | POST   | `/appointments/{appointmentId}/confirmation`          | Confirm or decline an appointment | Doctor         |
| 💤         | POST   | `/appointments/{appointmentId}/cancel`                | Reschedule an appointment         | Doctor/Patient |
| 💤         | POST   | `/appointments/{appointmentId}/feedback`              | Feedback an appointment           | Patient        |
| 💤         | GET    | `/appointments/{appointmentId}/feedback`              | Get feedbacks                     | Doctor/Patient |
| 💤         | GET    | `/appointments/{appointmentId}/feedback/{feedbackId}` | Get feedback by id                | Doctor/Patient |
| 💤         | POST   | `/files`                                              | Update files                      | Patient        |
| 💤         | GET    | `/files`                                              | Get all files                     | Patient        |
| 💤         | GET    | `/files/{fileId}`                                     | Get a file by id                  | Doctor/Patient |
| 💤         | POST   | `/files/{fileId}/access`                              | Create a file access              | Patient        |
| 💤         | GET    | `/files/{fileId}/access`                              | Get all file access               | Patient        |
| 💤         | GET    | `/files/{fileId}/access/{accessId}`                   | Get a file access by id           | Patient        |
| 💤         | PUT    | `/files/{fileId}/access/{accessId}`                   | Update a file access              | Patient        |
| 💤         | DELETE | `/files/{fileId}/access/{accessId}`                   | Delete a file access              | Patient        |
| 💤         | POST   | `/medical-reports`                                    | Create a medical report           | Doctor         |
| 💤         | GET    | `/medical-reports`                                    | Get all medical reports           | Doctor         |
| 💤         | GET    | `/medical-reports/{medicalReportId}`                  | Get a medical report by id        | Doctor         |
| 💤         | PUT    | `/medical-reports/{medicalReportId}`                  | Update a medical report           | Doctor         |
| 💤         | DELETE | `/medical-reports/{medicalReportId}`                  | Delete a medical report           | Doctor         |


# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.