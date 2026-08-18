import { useFieldContext } from "#/hooks/form.ts";
import { Field, FieldDescription, FieldError, FieldLabel } from "../ui/field";
import { Input } from "../ui/input";

type Props = {
	label: string;
	placeholder?: string;
	description?: string;
};
export const TextField = ({ label, placeholder, description }: Props) => {
	const { name, state, handleBlur, handleChange } = useFieldContext<string>();
	const isInvalid = state.meta.isTouched && !state.meta.isValid;

	return (
		<Field data-invalid={isInvalid}>
			<FieldLabel htmlFor={name}>{label}</FieldLabel>
			<Input
				type="text"
				id={name}
				name={name}
				value={state.value}
				onBlur={handleBlur}
				onChange={(e) => handleChange(e.target.value)}
				aria-invalid={isInvalid}
				placeholder={placeholder}
				autoComplete="off"
			/>
			{description && <FieldDescription>{description}</FieldDescription>}
			{isInvalid && <FieldError errors={state.meta.errors} />}
		</Field>
	);
};
