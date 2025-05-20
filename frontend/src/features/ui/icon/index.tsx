import { useState } from "react";
import { FaUserCircle } from "react-icons/fa";

type IconProps = {
  src: string;
  alt?: string;
  size?: number;
};

export const Icon = ({ src, alt = "avatar", size = 40 }: IconProps) => {
  const [imgError, setImgError] = useState(false);

  if (imgError) {
    return <FaUserCircle size={size} />;
  }

  return (
    <img
      src={src}
      alt={alt}
      width={size}
      height={size}
      onError={() => setImgError(true)}
      style={{ borderRadius: "50%" }}
    />
  );
};
