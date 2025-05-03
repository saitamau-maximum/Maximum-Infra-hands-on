import { MdSend } from "react-icons/md";

import styles from "./Button.module.css";

export const Button = () => {
  return(
    <button
      type="submit"
      className={styles.button}
    >
      <MdSend size={20} className={styles.icon} />
    </button>
  )
}