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


| Completed | Method | Endpoint                                                          | Description                              | User Role      |
| --------- | ------ | ----------------------------------------------------------------- | ---------------------------------------- | -------------- |
| ✅         | POST   | `/appointments`                                                   | Create an appointment via event          | Patient        |
| ✅         | GET    | `/appointments`                                                   | Get all appointments                     | Doctor/Patient |
| ✅         | GET    | `/appointments/{appointmentId}`                                   | Get an appointment by id                 | Doctor/Patient |
| ✅         | PUT    | `/appointments/{appointmentId}`                                   | Update an appointment                    | Patient        |
| ✅         | POST   | `/appointments/{appointmentId}/confirm`                           | Confirm or decline an appointment        | Doctor         |
| ✅         | POST   | `/appointments/{appointmentId}/cancel`                            | Cancel an appointment                    | Doctor/Patient |
| ✅         | POST   | `/appointments/{appointmentId}/feedbacks`                         | Add feedback to an appointment via event | Patient        |
| ✅         | GET    | `/appointments/{appointmentId}/feedbacks`                         | Get feedbacks                            | Doctor/Patient |
| ✅         | GET    | `/appointments/{appointmentId}/feedbacks/{feedbackId}`            | Get feedback by id                       | Doctor/Patient |
| ✅         | GET    | `/appointments/{appointmentId}/files`                             | Get all files attached to an appointment | Doctor         |
| ✅         | POST   | `/files`                                                          | Update files                             | Patient        |
| ✅         | GET    | `/files`                                                          | Get all files                            | Patient        |
| ✅         | GET    | `/files/{fileId}`                                                 | Get a file by id                         | Patient        |
| ✅         | POST   | `/files/{fileId}/access`                                          | Create a file access                     | Patient        |
| ✅         | GET    | `/files/{fileId}/access`                                          | Get all file access                      | Patient        |
| ✅         | POST   | `/appointments/{appointmentId}/medical-reports`                   | Create a medical report                  | Doctor         |
| ✅         | GET    | `/appointments/{appointmentId}/medical-reports`                   | Get all medical reports                  | Doctor         |
| ✅         | GET    | `/appointments/{appointmentId}/medical-reports/{medicalReportId}` | Get a medical report by id               | Doctor         |


# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.