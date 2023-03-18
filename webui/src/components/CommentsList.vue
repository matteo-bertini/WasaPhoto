<script>
export default {
	props: ["PhotoOwner","PhotoId",'CommentsList'],
	data(){
		return {
			Comments: []
		}
	},
	created(){
		this.Comments=this.CommentsList;
		return;
	},
	methods: {
		UpdateCommentList(CommentId){
			this.Comments = [];
			return;

		}
	}
	
}
</script>

<template>
	<div v-if="this.CommentsList.length == 0">
		<p>Ancora nessun commento</p>
		
	</div>
	<div v-else>
		<div v-for="Comment in CommentsList" :key ="Comment.CommentId">
			<Comment
			:PhotoOwner = "PhotoOwner"
			:PhotoId = "PhotoId"
			:CommentId = "Comment.CommentId" 
			:CommentAuthor="Comment.CommentAuthor" 
			:CommentText="Comment.CommentText"
			@commentdeleted = "UpdateCommentList"		
			>
		</Comment>
		</div>


	</div>
</template>

<style>
</style>
