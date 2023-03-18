<script>
export default {
	props: ["PhotoOwner","PhotoId","CommentId","CommentAuthor","CommentText"],
	emits: ["commentdeleted"],
	data() {
		return {
			isAuthor : false
		}
	},
	methods: {
		async deleteComment(){
			try {
				let config = {
					headers: {
						"Authorization": `Bearer ${localStorage.getItem("Authstring")}`
					}
				
				
				};
				this.$axios.delete("/users/"+this.PhotoOwner+"/photos/"+this.PhotoId+"/comments/"+this.CommentId,config);

				// Notifico al componente padre che ho eliminato la foto,così la tolga dalla lista 
				this.$emit("commentdeleted",this.CommentId);
				return;
			
			}catch(e){
				if(e.response.status == 500){
					this.$emit("Internal Server Error");
					return;


				}

			}

		}
	},
	mounted(){
		if(this.CommentAuthor !== `${localStorage.getItem("Username")}`){
			this.isAuthor=false;
			return;
		}
		else{
			this.isAuthor=true;
			return;

		}
	}
	
}
</script>

<template>
	<div style="display: flex; flex-direction: column;">
		<div style="display: flex; flex-direction: row; align-items: center;">
			<span style="font-size:larger;">
				{{CommentAuthor}}
			</span>
			<button class="btn" v-if="isAuthor" @click="deleteComment" id="deleteCommentButton" style="border: none;">
				<i class="fa-regular fa-trash-can"></i>
			</button>
		</div>
		<span>
			{{CommentText}}
		</span>
		
		


	</div>
</template>

<style>
	#deleteCommentButton:hover{
		color:#8b0000
	}

</style>
