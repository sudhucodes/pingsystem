import "dotenv/config";
import { Config } from "./types.js";

export const config: Config = {
    botToken: process.env.BOT_TOKEN ?? "",
    chatId: process.env.CHAT_ID ?? "",
};

if (!config.botToken || !config.chatId) {
    throw new Error("Missing BOT_TOKEN or CHAT_ID in .env");
}
