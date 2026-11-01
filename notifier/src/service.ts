// coolinary.notifier.v1.NotifierService
import { status } from "@grpc/grpc-js";
import { db } from "./db/client";
import { deliveries } from "./db/schema";
import type { SendEmailRequest, SendEmailResponse } from "./gen/coolinary/notifier/v1/notifier";
import { type NotifierServiceServer, NotifierServiceService } from "./gen/coolinary/notifier/v1/notifier";
import { type Registration, RpcError, unary } from "./grpc/server";
import { sendEmail } from "./mailer";

const EMAIL = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const implementation: NotifierServiceServer = {
  sendEmail: unary<SendEmailRequest, SendEmailResponse>(NotifierServiceService.sendEmail.path, async (req) => {
    if (!EMAIL.test(req.to)) throw new RpcError(status.INVALID_ARGUMENT, "to must be an e-mail address");
    if (!req.subject) throw new RpcError(status.INVALID_ARGUMENT, "subject is required");
    try {
      const messageId = await sendEmail(req.to, req.subject, req.textBody, req.htmlBody);
      await db
        .insert(deliveries)
        .values({ channel: "email", recipient: req.to, subject: req.subject, status: "sent", providerId: messageId });
      return { messageId };
    } catch (err) {
      await db
        .insert(deliveries)
        .values({ channel: "email", recipient: req.to, subject: req.subject, status: "failed", error: String(err) });
      throw new RpcError(status.UNAVAILABLE, `smtp: ${(err as Error).message}`);
    }
  }),
};

export const notifierService: Registration = {
  name: "coolinary.notifier.v1.NotifierService",
  definition: NotifierServiceService,
  implementation,
};
