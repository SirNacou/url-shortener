import {
	getApiUrlsQueryKey,
	postApiShortenMutation,
} from "#/api/gen/@tanstack/react-query.gen.ts";
import {
	type ShortenRequestBodyWritableZodInput,
	zShortenRequestBodyWritable,
} from "#/api/gen/zod.gen.ts";
import { SubmitButton } from "#/components/form/submit-button.tsx";
import { TextField } from "#/components/form/text-field.tsx";
import { UrlField } from "#/components/form/url-field.tsx";
import { Button } from "#/components/ui/button.tsx";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "#/components/ui/dialog.tsx";
import { FieldGroup } from "#/components/ui/field.tsx";
import { toast } from "#/components/ui/toast.tsx";
import { fieldContext, formContext } from "#/hooks/form.ts";
import { createFormHook } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { NumberField } from "../../components/form/number-field";

type Props = {};

const { useAppForm } = createFormHook({
	fieldComponents: {
		TextField,
		UrlField,
		NumberField,
	},
	formComponents: {
		SubmitButton,
	},
	fieldContext,
	formContext,
});

const defaultValues: ShortenRequestBodyWritableZodInput = {
	custom_code: "",
	expires_in_days: 30n,
	url: "",
};

export const CreateUrlDialog = (_props: Props) => {
	const queryClient = useQueryClient();
	const [open, setOpen] = useState(false);
	const { mutateAsync } = useMutation({
		...postApiShortenMutation(),
		onSuccess: () => {
			form.reset();
			setOpen(false);
			toast.add({ description: "URL created successfully", type: "success" });
			queryClient.invalidateQueries({ queryKey: getApiUrlsQueryKey() });
		},
		onError: () => {
			toast.add({ description: "Failed to create URL", type: "error" });
		},
	});

	const form = useAppForm({
		defaultValues,
		validators: {
			onChange: zShortenRequestBodyWritable,
			onSubmit: zShortenRequestBodyWritable,
		},
		onSubmit: async ({ value }) => {
			await mutateAsync({
				body: {
					expires_in_days: Number(value.expires_in_days),
					url: value.url,
					custom_code: value.custom_code,
				},
			});
		},
	});

	return (
		<Dialog open={open} onOpenChange={setOpen}>
			<DialogTrigger render={<Button variant="default">Create URL</Button>} />
			<DialogContent>
				<form
					className="flex flex-col gap-3"
					onSubmit={(e) => {
						e.preventDefault();
						e.stopPropagation();
						form.handleSubmit();
					}}
				>
					<DialogHeader>
						<DialogTitle className="text-center font-bold text-xl">
							Create URL
						</DialogTitle>
					</DialogHeader>
					<FieldGroup>
						<form.AppField name="custom_code">
							{(field) => (
								<field.TextField label="Code (optional)" placeholder="abcxyz" />
							)}
						</form.AppField>

						<form.AppField name="url">
							{(field) => (
								<field.UrlField label="URL" placeholder="https://example.com" />
							)}
						</form.AppField>

						<form.AppField name="expires_in_days">
							{(field) => (
								<field.NumberField label="Expires (days)" placeholder="30" />
							)}
						</form.AppField>
					</FieldGroup>
					<DialogFooter>
						<form.AppForm>
							<form.SubmitButton />
						</form.AppForm>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
};
