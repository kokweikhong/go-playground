import NextAuth from "next-auth/next";
import CredentialsProvider from "next-auth/providers/credentials";

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
        console.log(credentials);
        const res = await fetch("http://localhost:8080/api/v1/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(credentials),
        });

        const user = await res.json();
        console.log(user);
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
      console.log(account);
      // add 20 seconds
      return { ...token, ...user };
    },
    async session({ session, token, user }) {
      console.log("session", session, token);
      session.user = token as any;
      return session;
    },
  },
});

export { handler as GET, handler as POST };
