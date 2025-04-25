import { Form } from "../../ui/Form";
import { useForm } from "react-hook-form";
import { LoginFormData } from "../types/LoginFormDate";
import { loginUser } from "../api/login";
import styles from "./LoginPage.module.css";
import { useAuth } from "../hooks/useAuth";

export const LoginPage = () => {
  const {
    register,
    handleSubmit,
  } = useForm<LoginFormData>()
  const {user, loading} = useAuth();
  if (loading) return <div>Loading...</div>;
  return (
    <div className={styles.container}>
      <h1>Login</h1>
      <form className={styles.form} onSubmit={handleSubmit(loginUser)}>
        <Form.Field>
          <Form.Label label="Email" />
          <Form.Input
            type="email"
            id="email"
            required
            placeholder="Email"
            {...register("email")}
          />
          <Form.Label label="Password" />
          <Form.Input
            type="password"
            id="password"
            required
            placeholder="Password"
            {...register("password")}
          />
          {!user && (
            <Form.Button type="submit" >
              Login
            </Form.Button>
          )}
        </Form.Field>
      </form>
    </div>
  );
}