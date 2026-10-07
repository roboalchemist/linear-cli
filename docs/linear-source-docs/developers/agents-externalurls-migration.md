# Migrating to the externalUrls API

We’re upgrading the `AgentSession` update mutation to support multiple external URLs per session (from `externalLink` to `externalUrls`). 

Additionally, when a pull request is created by the agent, you should add that to the session’s `externalUrls` as well. This will unlock additional features related to pull requests in the future.

> [!NOTE]
> This change is backwards compatible, so for the moment, mutations updating `externalLink` will continue to work, but we highly recommend moving to the new `externalUrls` API so your agent can take advantage of new capabilities. 
> 
> **Writing to `externalLink` would overwrite all `externalUrls` set on the session.**

## Migration Prompt

Experimentally, we’ve prepared a short prompt that you can share with your coding agent to perform the API migration. You are free to modify this prompt to adapt to your own usage.

[migration.md →](https://gist.githubusercontent.com/itsmingjie/484129cf02f5855f3cc6fb7080e9ae21/raw/0392c230c0b7992892a1d6aa05b10133d3b219b1/agents-external-urls-migration.md)

## GraphQL

### Before

```graphql
mutation AgentSessionUpdateExternalUrl(
  $id: String!
  $input: AgentSessionUpdateExternalUrlInput!
) {
  agentSessionUpdateExternalUrl(id: $id, input: $input) {
    success
  }
}
```

### After

```graphql
mutation AgentSessionUpdate($id: String!, $input: AgentSessionUpdateInput!) {
  agentSessionUpdate(id: $id, input: $input) {
    success
  }
}

```

## TypeScript SDK (v70.0.0+)

### Before

```ts
import { LinearClient } from "@linear/sdk";

const linearClient = new LinearClient(...);
const session = await linearClient.agentSession(
  "<uuid>"
);

await session.update({
  externalLink: "https://example.com/agent/sessions/123",
});
// OR
await linearClient.agentSessionUpdateExternalUrl(
  "<uuid>",
  { externalLink: "https://example.com/agent/sessions/123" }
);
```

### After: Modifying external URLs

```ts
import { LinearClient } from "@linear/sdk";

const linearClient = new LinearClient(...);
const session = await linearClient.agentSession(
  "<uuid>"
);

await session.update({
  addedExternalUrls: [
    { url: "https://example.com/agent/sessions/123", label: "AgentName Web" },
  ],
  removedExternalUrls: ["https://example.com/old/sessions/123"],
  ],
});

```

### After: Adding a pull request

```ts
// ...

await session.update({
  addedExternalUrls: [{ url: "https://github.com/org/repo/pull/123", label: "Pull Request" }],
});
```