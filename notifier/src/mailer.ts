// SMTP delivery. Dev points at Mailpit (http://10.77.0.3:8025); prod at a real provider.
import nodemailer from "nodemailer";
import { config } from "./config";

const transport = nodemailer.createTransport(config.smtpUrl);

export async function sendEmail(to: string, subject: string, text: string, html?: string): Promise<string> {
  const info = await transport.sendMail({ from: config.mailFrom, to, subject, text, html: html || undefined });
  return info.messageId;
}
