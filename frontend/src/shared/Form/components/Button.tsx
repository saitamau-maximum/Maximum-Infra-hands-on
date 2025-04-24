import { ButtonHTMLAttributes } from "react"
import classNames from "classnames"
import styles from "./Button.module.css"

type ButtonProps = {
  variant?: "primary" | "secondary"
} & ButtonHTMLAttributes<HTMLButtonElement>

export const Button = ({
  variant = "primary",
  className,
  ...props
}: ButtonProps) => {
  return (
    <button
      className={classNames(styles.button, styles[variant], className)}
      {...props}
    />
  )
}
