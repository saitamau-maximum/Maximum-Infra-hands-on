import { Form } from "../../../shared/Form";
import styles from "./RegisterPage.module.css";

export const RegisterPage = () => {
  return (
    <div className={styles.container}>
      <h1>Register</h1>
      <form className={styles.form}>
        <Form.Field>
          <Form.Label label="Username" />
          <Form.Input
            type="text"
            id="username"
            name="username"
            required
            placeholder="Username"
          />
          <Form.Label label="Email" />
          <Form.Input
            type="email"
            id="email"
            name="email"
            required
            placeholder="Email"
          />
          <Form.Label label="Password" />
          <Form.Input
            type="password"
            id="password"
            name="password"
            required
            placeholder="Password"
          />
          <Form.Button type="submit" >
            Register
          </Form.Button>
        </Form.Field>
      </form>
    </div>
  );
}