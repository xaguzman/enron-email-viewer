export interface Email {
    Id: string,
    Bcc: string[],
    Body: string,
    Cc: string[],
    From: string,
    Subject: string,
    Date: string,
    To: string[]
}