import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";

export async function getUsers(): Promise<string> {
  const session = await getServerSession(authOptions);
  console.log(session?.user.accessToken);
  return session?.user.accessToken ?? "no token";
}
