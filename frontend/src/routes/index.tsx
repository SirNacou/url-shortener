import { getApiUrlsOptions } from "#/api/gen/@tanstack/react-query.gen.ts";
import CopyButton from "#/components/copy-button.tsx";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "#/components/ui/card.tsx";
import dayjs from "#/lib/dayjs.ts";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { CreateUrlDialog } from "./-urls/create-url-dialog";

export const Route = createFileRoute("/")({ component: Home });

function Home() {
	const {
		data: urls,
		isLoading,
		isError,
	} = useQuery({
		...getApiUrlsOptions(),
	});

	return (
		<div className="flex flex-col items-center p-8 gap-3">
			<div className="text-3xl font-bold">URL Shortener</div>
			<CreateUrlDialog />
			<div className="grid grid-cols-4">
				{urls?.items?.map((url) => (
					<Card key={url.code}>
						<CardHeader>
							<CardTitle>
								<CopyButton
									code={url.code}
									content={url.short_url}
									duration={1500}
								/>
							</CardTitle>
						</CardHeader>
						<CardContent className="flex flex-col gap-3">
							<p>URL: {url.target_url}</p>
							<p>Expires At: {dayjs(url.expires_at).format("YYYY-MM-DD")}</p>
						</CardContent>
					</Card>
				))}
			</div>
		</div>
	);
}
