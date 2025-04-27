import { LoginForm } from "../components/LoginForm";

import styles from "./LoginPage.module.css";

export const LoginPage = () => {
  return (
    <div className={styles.container}>
      <h1>Login</h1>
      <LoginForm />
    </div>
  );
}