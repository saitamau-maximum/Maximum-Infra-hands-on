import { useForm } from "react-hook-form";
import { Form } from "../../ui/Form";
import { uploadImageApi } from "../api";

type FormValues = {
  image: FileList;
};

export const ImageUploadPage = () => {
  const { register, handleSubmit, reset } = useForm<FormValues>();

  const onSubmit = async (data: FormValues) => {
    const file = data.image[0];
    if (!file) return;

    try {
      await uploadImageApi(file);
      alert("Upload successful!");
      reset(); // フォームをリセット
    } catch (err) {
      alert("Upload failed");
      console.error(err);
    }
  };

  return (
    <div className="p-4 border rounded w-80">
      <h1 className="text-lg mb-2">Upload Image</h1>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Form.Field>
          <Form.Label label="Upload Image" />
          <Form.Input
            type="file"
            accept="image/jpeg, image/png"
            {...register("image")}
          />
        </Form.Field>
        <button type="submit" className="mt-2 px-4 py-1 bg-blue-500 text-white rounded">
          Upload
        </button>
      </form>
    </div>
  );
};
