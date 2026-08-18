import { useFormContext } from "#/hooks/form.ts";
import { Button } from "../ui/button";

type Props = {};
export const SubmitButton = (_: Props) => {
	const form = useFormContext();
	return (
		<form.Subscribe
			selector={(state) => ({
				canSubmit: state.canSubmit,
				isSubmitting: state.isSubmitting,
			})}
		>
			{({ canSubmit, isSubmitting }) => (
				<Button type="submit" disabled={!canSubmit || isSubmitting}>
					{isSubmitting ? "Submitting..." : "Submit"}
				</Button>
			)}
		</form.Subscribe>
	);
};
