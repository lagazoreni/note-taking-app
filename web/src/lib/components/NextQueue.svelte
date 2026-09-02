<script lang="ts">
	import type { Question } from '$lib/types/question';

	export let sections: { id: string; title: string; items: Question[] }[] = [];
</script>

{#if sections.length === 0}
	<p>Nothing in Next. Capture a question from a note, or set a due date.</p>
{:else}
	<div class="next-queue">
		{#each sections.filter((section) => section.items.length > 0) as section (section.id)}
			<section aria-labelledby={`next-section-${section.id}`}>
				<h2 id={`next-section-${section.id}`}>{section.title}</h2>
				<ul>
					{#each section.items as question (question.id)}
						<li>
							<a
								href={question.linkedNotes[0]
									? `/questions/${question.id}?noteId=${question.linkedNotes[0].id}`
									: `/questions/${question.id}`}
							>
								{question.questionText}
							</a>
						</li>
					{/each}
				</ul>
			</section>
		{/each}
	</div>
{/if}
