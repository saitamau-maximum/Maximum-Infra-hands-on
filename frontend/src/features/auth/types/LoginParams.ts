import { LoginFormData } from "./LoginFormDate";

export type LoginParams = {
  data: LoginFormData;
  refetch: () => void;
};