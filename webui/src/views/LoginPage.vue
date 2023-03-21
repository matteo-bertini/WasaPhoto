<script>
export default {
	// Dichiarazione del reactive state
	data: function () {
		return {
			errormsg: null,
			loading: false,
			checked : false,
			Username : "",
			
		}
	},
	methods: {
		doLogin :async function () {
			this.loading = true;
    		this.errormsg = null;
			try{
				let response1 = await this.$axios.post("/session", {Username: this.Username});
            	localStorage.setItem('Authstring',response1.data.Identifier);
           		localStorage.setItem('Username',this.Username);
				if(this.checked == false){
					let response2 = await this.$axios.post("/users/", {Username: this.Username},{headers: { Authorization: `Bearer ${localStorage.getItem("Authstring")}`}});
            		this.$router.push("/users/"+this.Username);
					this.loading = false;
            		return
				}
				else{
					this.$router.push("/users/"+this.Username)
					this.loading = false;
            		return

				}
			}catch(e){
				if(e.response.status == 403){
					this.errormsg = "Oops! Sembra che tu abbia già un profilo esistente. Per favore spunta la casella in basso a sinistra nel pannello di login.";
					this.loading = false;
					return;
				}else{
					this.errormsg = e.toString();
					this.loading = false;
					return;

				}
				

			}
		}
	}
	
}


</script>

<template>
	<div class = "container-fluid" >
		<div class="row">
			<div class="col">
					<ErrorMsg v-if="errormsg" :Message="errormsg"></ErrorMsg>
			</div>
		</div>
		<div class="row" style="display: flex; justify-content: center; align-items: center; vertical-align: middle; margin-top: 250px;" >
			<div class="col">
				<div class="card" style="background-color:rgba(255, 255, 255, 0.285);">
					<p style="text-align: center; color: black; font-size:50px;">Welcome to WasaPhoto!</p>
					<p style="text-align: center; color: black; font-size:18px">
						Keep in touch with your friends by sharing photos of special moments, thanks to WASAPhoto! <br>
						You can upload your photos directly from your PC, and they will be visible to everyone following you.
					</p>	

				</div>
			</div>
			<div class="col">
				<div class="card" style="background-color:rgba(255, 255, 255, 0.285); display: flex; justify-content: center; align-items: center; text-align: center;">
					<div>
						<h1 style="font-family:Arial;">LOGIN</h1>
					</div>
					<div>
						<input name="LoginTextArea" class="form-control" v-model="Username" placeholder="Enter here your username" style="width:fit-content" >
						<button type="button" @click="doLogin" class="btn btn-dark" :disabled = "Username==null || Username.length<3 || Username.length>30 || Username.trim().length<3"  style="width:fit-content; margin-top: 10px;">Login/Register</button>
					</div>
					
					<div style="display: inline-flex; align-self: flex-start; margin-left: 15px; margin-bottom: 5px;">
						<input type="checkbox" id="checkbox" v-model="checked"/>
						<label for="checkbox" style="margin-left: 5px;">I have a profile</label>

					</div>
					
						
				</div>


			</div>

		</div>	
	</div>
</template>

<style>
	


</style>
