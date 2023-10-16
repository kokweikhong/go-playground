import NextAuth from "next-auth";

declare module "next-auth" {
  interface Session {
    user: {
      id: number;
      username: string;
      email: string;
      role?: string;
      accessToken: string;
      maxAge: number;
    };
  }
}
