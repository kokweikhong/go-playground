// import { NextRequest } from "next/server";
import { NextRequestWithAuth } from "next-auth/middleware";

export function isAuthenticated(req: NextRequestWithAuth): boolean {
  const cookies = req.nextauth.token;
  console.log("from isAuthenticated");

  console.log(cookies?.error);
  return true;
}
