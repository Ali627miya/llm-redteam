/**
 * Vercel AI SDK: the redteam CLI tests HTTP. Point `target.url` at your route
 * (e.g. POST /api/chat) and set `body_template` / `response_path` to match
 * the JSON your handler returns (often `{ "text": "..." }` or stream off in v2 — use a non-streaming test route for scans).
 */
// app/api/redteam-echo/route.ts — minimal non-streaming target for scans:
// export async function POST(req: Request) {
//   const { messages } = await req.json();
//   const result = await streamText({ model: openai("gpt-4o-mini"), messages });
//   return Response.json({ text: await result.text });
// }

export {};
