"use client";

import React from "react";
import { useSession } from "next-auth/react";
import { User } from "@/types/user";

export default function Page() {
  const { data: session } = useSession();
  console.log(session);

  const [users, setUsers] = React.useState<User[]>([]);

  async function helloWord() {
    const res = await fetch("http://localhost:8080/api/v1");
    const data = await res.json();

    console.log(data);
  }

  async function getUsers() {
    const res = await fetch("http://localhost:8080/api/v1/users", {
      headers: {
        Authorization: `Bearer ${session?.user.accessToken}`,
      },
    });
    const data = await res.json();

    setUsers(data.data);
  }

  return (
    <main>
      <h1>Users Page</h1>
      <p>This is a protected page</p>
      <button onClick={helloWord}>Get Hello World</button>
      <div>
        <button onClick={getUsers}>Get users</button>
        <ul>
          {users?.map((user) => (
            <li key={user.id}>{user.email}</li>
          ))}
        </ul>
      </div>
    </main>
  );
}
