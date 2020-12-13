import React, {useState} from 'react';
import Router from 'next/router';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  function handleSubmit(e) {
    e.preventDefault();
    //call api
    fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        'username': username,
        'password': password,
      }),
    })
      .then((r) => {
        return r.json();
      })
      .then((data) => {
        Router.push(`/user`);
      });
  }
  return (
    <form onSubmit={handleSubmit}>
      <p>ログイン</p>
      <input
        name="username"
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <input
        name="password"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <input type="submit" value="Submit" />
    </form>
  );
};

export default Login;
