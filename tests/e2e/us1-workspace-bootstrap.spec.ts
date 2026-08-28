import { test, expect, type Page } from "@playwright/test";

const externalNetworkBlock = async (page: Page) => {
  await page.route("**/*", async (route) => {
    const url = new URL(route.request().url());
    if (url.hostname !== "127.0.0.1" && url.hostname !== "localhost")
      return route.abort();
    return route.continue();
  });
};

test("bootstraps a workspace before enabling note capture on a fresh database", async ({
  page,
}) => {
  await externalNetworkBlock(page);

  let noteRequests = 0;
  await page.route("**/api/v1/notes**", async (route) => {
    noteRequests += 1;
    return route.continue();
  });

  await page.goto("/");
  await expect(
    page.getByRole("heading", { name: /create your first workspace/i }),
  ).toBeVisible();
  await expect(
    page.getByRole("link", { name: /create workspace/i }),
  ).toBeVisible();

  // Direct navigation is guarded too: it must not reach the editor or issue a note request.
  await page.goto("/notes/new");
  await expect(
    page.getByText(/select or create a workspace first/i),
  ).toBeVisible();
  expect(noteRequests).toBe(0);

  await page.getByRole("link", { name: /go to workspaces/i }).click();
  await page.getByLabel("New workspace").fill("Bootstrap workspace");
  await page.getByRole("button", { name: "Create" }).click();

  await expect(
    page.getByRole("option", { name: "Bootstrap workspace" }),
  ).toBeVisible();
  await expect(page.getByLabel("Current workspace")).toHaveValue(/.+/);
  await expect(page.getByText(/workspace created/i)).toBeVisible();

  await page.getByRole("link", { name: "Notes" }).click();
  await page.getByRole("link", { name: "New note" }).click();
  await expect(
    page.getByRole("heading", { name: "Create a note" }),
  ).toBeVisible();
  expect(noteRequests).toBe(0);
});
