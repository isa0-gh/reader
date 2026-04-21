<script lang="ts">
  import { onMount } from 'svelte';

  let view = 'login';
  let email = '';
  let password = '';
  let name = '';
  let error = '';
  let user: any = null;
  let token = '';

  const API_URL = '/api/v1';

  onMount(() => {
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    if (savedToken && savedUser) {
      token = savedToken;
      user = JSON.parse(savedUser);
      view = 'home';
    }
  });

  async function handleLogin() {
    error = '';
    try {
      const res = await fetch(`${API_URL}/users/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data || 'Login failed');
      
      token = data.token;
      user = data.user;
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      view = 'home';
    } catch (e: any) {
      error = e.message;
    }
  }

  async function handleRegister() {
    error = '';
    try {
      const res = await fetch(`${API_URL}/users/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, name })
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data || 'Registration failed');
      
      // Auto login or switch to login
      view = 'login';
      password = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  function handleLogout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    token = '';
    user = null;
    view = 'login';
  }
</script>

<main>
  {#if view === 'login'}
    <h1>Login</h1>
    {#if error}<p class="error">{error}</p>{/if}
    <input type="email" placeholder="Email" bind:value={email} />
    <input type="password" placeholder="Password" bind:value={password} />
    <button on:click={handleLogin}>Login</button>
    <button class="link-btn" on:click={() => view = 'register'}>Don't have an account? Register</button>

  {:else if view === 'register'}
    <h1>Register</h1>
    {#if error}<p class="error">{error}</p>{/if}
    <input type="text" placeholder="Name" bind:value={name} />
    <input type="email" placeholder="Email" bind:value={email} />
    <input type="password" placeholder="Password" bind:value={password} />
    <button on:click={handleRegister}>Register</button>
    <button class="link-btn" on:click={() => view = 'login'}>Already have an account? Login</button>

  {:else if view === 'home'}
    <h1>Welcome, {user?.name || 'User'}</h1>
    <p>Logged in as: {user?.email}</p>
    <button style="margin-top: 2rem;" on:click={handleLogout}>Logout</button>
  {/if}
</main>
