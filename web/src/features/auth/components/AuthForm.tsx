// Note: separate this if sign in and up diverge
import { Form, FormInputField, FormSubmitButton, FormTitle } from "@/components";
import { useForm } from "react-hook-form";

type AuthFormProps = {
    onSubmit: (values: { email: string; password: string }) => void;
    isSubmitting?: boolean;
    formType: "signin" | "signup";
}

export default function AuthForm({
    onSubmit,
    isSubmitting = false,
    formType,
}: AuthFormProps) {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<{ email: string; password: string }>();
    return (
        <Form onSubmit={handleSubmit(onSubmit)}>
            <FormTitle text={formType === "signin" ? "Sign In" : "Sign Up"} />
            <FormInputField  
                label="Email"
                id="email"
                register={register("email", { 
                    required: "Email is required" 
                })}
                error={errors.email}
                placeholder="Email"
            />
            {/* TODO: Add toggleable visibility text input */}
            <FormInputField  
                label="Password"
                id="password"
                register={register("password", { 
                    required: "Password is required"
                })}
                error={errors.password}
                placeholder="Password"
            />      
            <FormSubmitButton disabled={isSubmitting}>
                {formType === "signin" ? "Sign In" : "Sign Up"}
            </FormSubmitButton>
        </Form>
    );
}