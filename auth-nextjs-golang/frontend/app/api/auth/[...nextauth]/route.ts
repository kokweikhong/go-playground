import NextAuth from "next-auth/next";
import CredentialsProvider from "next-auth/providers/credentials";
import { JWT } from "next-auth/jwt";

// refresh token if access token expired
async function refreshToken(token: JWT) {
  console.log("refresh token");
  const res = await fetch("http://localhost:8080/api/v1/auth/refresh", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username: token.username,
      refreshToken: token.refreshToken,
    }),
  });

  console.log(res.ok);

  if (!res.ok) {
    return {
      ...token,
      error: "RefreshTokenExpired",
    };
  }

  const newAccessToken = await res.json();
  const newAccessTokenExpiry = newAccessToken.accessTokenExpiry;

  return {
    ...token,
    accessToken: newAccessToken.accessToken,
    accessTokenExpiry: newAccessTokenExpiry,
  };
}

const handler = NextAuth({
  // session: {
  //   strategy: "jwt",
  //   // 10 seconds
  //   maxAge: 10,
  // },
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        username: { label: "Username", type: "text" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials, req) {
        // console.log(credentials);
        const res = await fetch("http://localhost:8080/api/v1/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(credentials),
        });

        const user = await res.json();
        // console.log(user);
        // const user = { username: "J Smith", password: "1234" };
        // const user = { id: 1, name: "J Smith", email: "" };

        if (user) {
          return user;
        } else {
          return null;
        }
      },
    }),
  ],
  pages: {
    signIn: "/auth/login",
  },
  secret: process.env.NEXTAUTH_SECRET,
  callbacks: {
    async jwt({ token, user, account }) {
      // console.log(account);
      // console.log("token");
      // console.log(token);
      // console.log("user");
      // console.log(user);
      // return { ...token, ...user };
      token = { ...token, ...user };

      // convert access token expiry number to Date
      // token.accessTokenExpiry = new Date(token.accessTokenExpiry);
      console.log(Date.now(), token.accessTokenExpiry);
      // Date.now() sub 30 seconds
      // Date.now() - 30000
      if (
        Date.now() - 5000 <
        (token.accessTokenExpiry as unknown as number) * 1000
      ) {
        return token;
      } else {
        const newToken = await refreshToken(token);
        console.log("new token");
        console.log(newToken);

        return { ...newToken };
      }
    },
    async session({ session, token, user }) {
      // console.log("session", session, token);
      session.user = token as any;
      return session;
    },
  },
});

export { handler as GET, handler as POST };
