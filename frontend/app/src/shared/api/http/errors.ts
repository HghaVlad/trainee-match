export class AppError extends Error {
  status?: number
  constructor(message: string, status?: number) {
    super(message)
    this.name = 'AppError'
    this.status = status
  }
}

export class SessionExpiredError extends AppError {
  constructor(message = 'Session expired') {
    super(message, 401)
    this.name = 'SessionExpiredError'
  }
}
